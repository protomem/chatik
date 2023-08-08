interface BaseEntity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserEntity extends BaseEntity {
  nickname: string;
  email: string;
  isVerified: boolean;
}

export interface ChannelEntity extends BaseEntity {
  title: string;
  user: UserEntity;
}

export interface MessageEntity extends BaseEntity {
  content: string;
  channelId: string;
  user: UserEntity;
}
