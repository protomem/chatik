import { RootState } from "../store";

const selectAuthModule = (state: RootState) => state.auth;

export const selectToken = (state: RootState) => selectAuthModule(state).token;
export const selectUser = (state: RootState) => selectAuthModule(state).user;
export const selectIsAuth = (state: RootState) =>
  selectToken(state) !== "" && selectUser(state) !== null;
