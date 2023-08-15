import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { UserEntity } from "../../entities/entities";
import { authStorage } from "./auth.storage";

interface AuthState {
  currentUser: UserEntity | null;
  accessToken: string;
}

const initialState: AuthState = {
  currentUser: authStorage.getCurrentUser(),
  accessToken: authStorage.getAccessToken(),
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }: PayloadAction<AuthState>) => {
      if (!payload.currentUser || !payload.accessToken) {
        return;
      }

      state.currentUser = payload.currentUser;
      authStorage.setCurrentUser(payload.currentUser);

      state.accessToken = payload.accessToken;
      authStorage.setAccessToken(payload.accessToken);
    },

    clearCredentials: (state) => {
      state.currentUser = null;
      authStorage.removeCurrentUser();

      state.accessToken = "";
      authStorage.removeAccessToken();
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
