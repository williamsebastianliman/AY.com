import api from "../lib/api";

export interface ThreadMedia {
  id: string;
  image_url: string;
  extension: string;
}
export interface RawThread {
  id: string;
  user_id: string;
  community_id?: string;
  content: string;
  like_count: number;
  comment_count: number;
  share_count: number;
  view_count: number;
  created_at: {
    year: number;
    month: number;
    day: number;
    hour: number;
    minute: number;
    second: number;
    timezone: string;
  };
  updated_at: {
    year: number;
    month: number;
    day: number;
    hour: number;
    minute: number;
    second: number;
    timezone: string;
  };
  is_premium: boolean;
  repost_id: string;
  is_pinned: boolean;
}

export interface Thread {
  id: string;
  user_id: string;
  author: string;
  username: string;
  avatar: string;
  content: string;
  community_id?: string;
  like_count: number;
  repost_count: number;
  is_reposted: boolean;
  comment_count: number;
  share_count: number;
  view_count: number;
  created_at: string;
  updated_at: string;
  media: ThreadMedia[];
  liked_by_me?: boolean;
  bookmarked: boolean;
  is_premium: boolean;
  repost_id: string;
  repost_thread: Thread | null;
  is_pinned: boolean;
}

export async function fetchThreadMedia(
  threadId: string
): Promise<ThreadMedia[]> {
  try {
    const res = await api.get<ThreadMedia[]>(`threads/${threadId}/media`);
    return res.data;
  } catch {
    return [];
  }
}

export async function fetchLiked(
  threadId: string,
  userId: string
): Promise<boolean> {
  try {
    console.log("uidl: ", userId);
    const res = await api.post<{ liked: boolean }>(
      `/threads/${threadId}/liked`,
      { user_id: userId },
      { withCredentials: true }
    );
    return res.data.liked;
  } catch {
    return false;
  }
}

export async function fetchReposted(
  threadId: string,
  userId: string
): Promise<boolean> {
  try {
    const res = await api.post<{ reposted: boolean }>(
      `/threads/${threadId}/reposted`,
      { user_id: userId },
      { withCredentials: true }
    );
    return res.data.reposted;
  } catch {
    return false;
  }
}

export async function fetchRepostCount(threadId: string): Promise<number> {
  try {
    const res = await api.get<{ count: number }>(
      `/threads/${threadId}/repost/count`
    );
    return res.data.count;
  } catch {
    return 0;
  }
}

export async function fetchBookmarked(
  threadId: string,
  userId: string
): Promise<boolean> {
  try {
    const res = await api.post<{ bookmarked: boolean }>(
      `/threads/${threadId}/bookmarked`,
      { user_id: userId },
      { withCredentials: true }
    );
    return res.data.bookmarked;
  } catch {
    return false;
  }
}

export async function fetchLikeCount(threadId: string): Promise<number> {
  try {
    const res = await api.get<{ count: number }>(`/threads/${threadId}/likes`);
    return res.data.count;
  } catch {
    return 0;
  }
}

export async function fetchRepliesCount(threadId: string): Promise<number> {
  try {
    const res = await api.get<{ count: number }>(
      `/threads/${threadId}/replies/count`
    );
    return res.data.count;
  } catch {
    return 0;
  }
}

export async function fetchUserAndAvatar(userId: string): Promise<{
  author: string;
  username: string;
  avatar: string;
  is_premium: boolean;
}> {
  const userRes = await api.get<{
    id: string;
    name: string;
    username: string;
    profile_picture_id: string;
    is_premium: boolean;
  }>(`/user/${userId}`);
  const { name, username, profile_picture_id, is_premium } = userRes.data;
  console.log("is premium: ", is_premium);

  let avatar = "";
  if (profile_picture_id) {
    const mediaRes = await api.post<{ public_url: string }>(
      "/media/get-media",
      { id: profile_picture_id }
    );
    avatar = mediaRes.data.public_url;
  }
  return {
    author: name,
    username,
    avatar,
    is_premium: is_premium,
  };
}

export async function getThreadByID(id: string): Promise<RawThread> {
  const res = await api.get<RawThread>(`/threads/${id}`);
  return res.data;
}

export async function enrichThread(
  raw: RawThread,
  userId: string
): Promise<Thread> {
  const { author, username, avatar, is_premium } = await fetchUserAndAvatar(
    raw.user_id
  );

  const { year, month, day, hour, minute, second } = raw.created_at;
  const utcMs = Date.UTC(year, month - 1, day, hour, minute, second);
  const dt = new Date(utcMs);
  const createdAt = dt.toISOString();

  const media = await fetchThreadMedia(raw.id);
  const like_count = await fetchLikeCount(raw.id);
  const liked_by_me = await fetchLiked(raw.id, userId);
  const bookmarked = await fetchBookmarked(raw.id, userId);
  const comment_count = await fetchRepliesCount(raw.id);
  const repost_count = await fetchRepostCount(raw.id);
  const is_reposted = await fetchReposted(raw.id, userId);

  let repostThread: Thread | null = null;
  if (raw.repost_id && raw.repost_id !== "") {
    try {
      let originalRaw = await getThreadByID(raw.repost_id);
      while (originalRaw.repost_id && originalRaw.repost_id !== "") {
        originalRaw = await getThreadByID(originalRaw.repost_id);
      }
      repostThread = await enrichThread(originalRaw, userId);
    } catch (err) {
      repostThread = null;
      console.error(err as string);
    }
  }

  return {
    id: raw.id,
    user_id: raw.user_id,
    author,
    username,
    avatar,
    content: raw.content,
    community_id: raw.community_id,
    like_count,
    liked_by_me,
    comment_count,
    share_count: raw.share_count,
    view_count: raw.view_count,
    created_at: createdAt,
    updated_at: createdAt,
    media,
    bookmarked,
    is_premium,
    repost_id: raw.repost_id,
    repost_thread: repostThread,
    repost_count,
    is_reposted,
    is_pinned: raw.is_pinned,
  };
}
