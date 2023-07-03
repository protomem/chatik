import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { authActions } from "../store/auth/auth.slice";

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/auth`,
  }),
  endpoints: (builder) => ({
    login: builder.mutation({
      query: (credentials) => ({
        url: "/login",
        method: "POST",
        body: { ...credentials },
      }),
    }),
  }),
});

export const { useLoginMutation } = authApi;
