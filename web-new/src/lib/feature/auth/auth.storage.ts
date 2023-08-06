import { UserEntity } from "../../entities/entities";

export const authStorage = {
  getCurrentUser: (): UserEntity | null => {
    const user = localStorage.getItem("currentUser");
    if (user) {
      return JSON.parse(user) as UserEntity;
    }

    return null;
  },

  setCurrentUser: (user: UserEntity) => {
    localStorage.setItem("currentUser", JSON.stringify(user));
  },

  removeCurrentUser: (): void => {
    localStorage.removeItem("currentUser");
  },

  getAccessToken: (): string => {
    return localStorage.getItem("accessToken") || "";
  },

  setAccessToken: (token: string) => {
    localStorage.setItem("accessToken", token);
  },

  removeAccessToken: (): void => {
    localStorage.removeItem("accessToken");
  },
};
