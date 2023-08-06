import { RootState } from "../store";

const selectAuthModule = (state: RootState) => state.auth;

export const selectCurrentUser = (state: RootState) =>
  selectAuthModule(state).currentUser;
export const selectAccessToken = (state: RootState) =>
  selectAuthModule(state).accessToken;
export const selectIsAuthenticated = (state: RootState) =>
  selectCurrentUser(state) !== null && selectAccessToken(state) !== "";
