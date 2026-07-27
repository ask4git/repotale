export type Post = {
  slug: string;
  title: string;
  excerpt: string;
  tags: string[];
  prNumber: number;
  prUrl: string;
  publishedAt: string;
};

export type Repo = {
  owner: string;
  name: string;
  private: boolean;
};
