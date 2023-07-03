import { createSlice } from "@reduxjs/toolkit";

// TODO: recovery from localStorage
const initialState = {
  currentUser: null,
  accessToken: "",
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload }) => {
      state.currentUser = payload.user;
      state.accessToken = payload.accessToken;
    },

    logout: (state) => {
      state.currentUser = null;
      state.accessToken = "";
    },
  },
});

export const authReducer = authSlice.reducer;
export const authActions = authSlice.actions;
