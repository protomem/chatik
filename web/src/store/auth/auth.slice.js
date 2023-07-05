import { createSlice } from "@reduxjs/toolkit";

const getCurrentStoreFromStorage = () => {
  const user = JSON.parse(localStorage.getItem("user"));
  if (user) {
    return user;
  }

  return null;
};

const saveCurrentUserToStorage = (user) => {
  localStorage.setItem("user", JSON.stringify(user));
};

const removeCurrentUserFromStorage = () => {
  localStorage.removeItem("user");
};

const getAccessTokenFromStorage = () => {
  const accessToken = localStorage.getItem("accessToken");
  if (accessToken) {
    return accessToken;
  }

  return "";
};

const saveAccessTokenToStorage = (accessToken) => {
  localStorage.setItem("accessToken", accessToken);
};

const removeAccessTokenFromStorage = () => {
  localStorage.removeItem("accessToken");
};

// TODO: recovery from localStorage
const initialState = {
  currentUser: getCurrentStoreFromStorage(),
  accessToken: getAccessTokenFromStorage(),
  isLoggedIn:
    getCurrentStoreFromStorage() === null || getAccessTokenFromStorage() === ""
      ? false
      : true,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }) => {
      state.currentUser = payload.user;
      saveCurrentUserToStorage(payload.user);

      state.accessToken = payload.accessToken;
      saveAccessTokenToStorage(payload.accessToken);

      state.isLoggedIn = true;
    },

    logout: (state) => {
      state.currentUser = null;
      removeCurrentUserFromStorage();

      state.accessToken = "";
      removeAccessTokenFromStorage();

      state.isLoggedIn = false;
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
