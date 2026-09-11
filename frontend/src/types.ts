export interface User {
  id: number;
  email: string;
  username: string;
  is_verified: boolean;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
}

export interface Post {
  id: number;
  user_id: number;
  username: string;
  title: string;
  content: string;
  created_at: string;
  categories: Category[];
  likes: number;
  dislikes: number;
  user_vote: number; // 0, 1, -1
  comment_count: number;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  username: string;
  content: string;
  created_at: string;
  likes: number;
  dislikes: number;
  user_vote: number;
}

export interface VoteResponse {
  likes: number;
  dislikes: number;
  user_vote: number;
}

export interface TokenResponse {
  access_token: string;
  user: User;
}

export interface ApiError {
  error: string;
  code?: string;
}
