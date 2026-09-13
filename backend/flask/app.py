import pickle
from flask import Flask, request, jsonify

import nltk
from nltk.tokenize import word_tokenize
from nltk.corpus import stopwords, wordnet
from nltk.stem import WordNetLemmatizer, SnowballStemmer
from nltk import pos_tag
import string

nltk.download("punkt")
nltk.download("stopwords")
nltk.download("averaged_perceptron_tagger")
nltk.download("wordnet")

stop_words = set(stopwords.words("english"))
punctuation = set(string.punctuation)
lemmatizer = WordNetLemmatizer()
stemmer = SnowballStemmer("english")

def get_tag(treebank_tag: str):
    treebank_tag = treebank_tag.lower()
    if treebank_tag.startswith("j"):
        return wordnet.ADJ
    elif treebank_tag.startswith("v"):
        return wordnet.VERB
    elif treebank_tag.startswith("n"):
        return wordnet.NOUN
    elif treebank_tag.startswith("r"):
        return wordnet.ADV
    else:
        return wordnet.NOUN

def pre_process(sentence: str):
    token_list = []
    sentence = sentence.lower()
    tokens = word_tokenize(sentence)
    for token in tokens:
        if token not in stop_words and token not in punctuation and token.isalpha():
            token_list.append(token)
    lemmatized_token_list = []
    pos_tagged = pos_tag(token_list)
    for word, tag in pos_tagged:
        lemma = lemmatizer.lemmatize(word, get_tag(tag))
        lemmatized_token_list.append(lemma)

    return lemmatized_token_list

def preprocess_and_join(texts):
    return [" ".join(pre_process(t)) for t in texts]

with open("svm_full_pipeline.pkl", "rb") as f:
    model = pickle.load(f)

app = Flask(__name__)

@app.route("/predict", methods=["POST"])
def predict():
    payload = request.get_json()
    text = payload.get("text", "")
    pred = model.predict([text])[0]
    return jsonify({"prediction": int(pred)})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True)