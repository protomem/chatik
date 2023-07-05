const selectAuthModule = (state) => state.auth;

export const selectCurrentUser = (state) => selectAuthModule(state).currentUser;
export const selectAccessToken = (state) => selectAuthModule(state).accessToken;
export const selectIsLoggedIn = (state) => selectAuthModule(state).isLoggedIn;
