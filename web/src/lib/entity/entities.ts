export interface IUser {
  id: string;
  createdAt: Date;
  updatedAt: Date;
  nickname: string;
  email: string;
  isVerified: boolean;
}

export interface IChannel {
  id: string;
  createdAt: Date;
  updatedAt: Date;
  title: string;
  user: IUser;
}

export interface IMessage {
  id: string;
  createdAt: Date;
  updatedAt: Date;
  content: string;
  channelId: string;
  user: IUser;
}
