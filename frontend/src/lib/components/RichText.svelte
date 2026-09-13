<script lang="ts">
    export let content: string = "";
    const pattern = /(@[a-zA-Z0-9_]+|#[a-zA-Z0-9_]+)/g;
    function parseRichText(text: string) {
      const result: { text: string; isTag: boolean; key: string }[] = [];
      let lastIndex = 0;
      let match: RegExpExecArray | null;
      let i = 0;
      while ((match = pattern.exec(text)) !== null) {
        if (match.index > lastIndex) {
          result.push({ text: text.slice(lastIndex, match.index), isTag: false, key: `text-${i++}` });
        }
        result.push({ text: match[0], isTag: true, key: `token-${i++}` });
        lastIndex = pattern.lastIndex;
      }
      if (lastIndex < text.length) {
        result.push({ text: text.slice(lastIndex), isTag: false, key: `text-${i++}` });
      }
      return result;
    }
    $: tokens = parseRichText(content);
  </script>
  <span class="text">
    {#each tokens as token (token.key)}
      {#if token.isTag}
        <span class="token">{token.text}</span>
      {:else}
        {token.text}
      {/if}
    {/each}
  </span>
  <style>
    .text {
      margin: 0.75rem 0;
      line-height: 1.5;
      white-space: pre-wrap;
    }
    .token {
      color: #228be6;
      font-weight: 500;
    }
  </style>