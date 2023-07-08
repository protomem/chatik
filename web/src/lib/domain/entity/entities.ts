export interface IUser {
  id: number;
  createdAt: Date;
  updatedAt: Date;
  name: string;
  email: string;
  isVerified: boolean;
}
