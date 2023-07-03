const selectAuthModule = (state) => state.auth;

export const selectCurrentUser = (state) => selectAuthModule(state).user;
export const selectAccessToken = (state) => selectAuthModule(state).token;
