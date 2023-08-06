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
