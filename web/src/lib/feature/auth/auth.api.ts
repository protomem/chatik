import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { UserEntity } from "../../entities/entities";

interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  accessToken: string;
  user: UserEntity;
}

interface RegisterRequest {
  nickname: string;
  email: string;
  password: string;
}

interface RegisterResponse {
  accessToken: string;
  user: UserEntity;
}

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/auth`,
  }),
  endpoints: (builder) => ({
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: ({ email, password }) => ({
        url: "/login",
        method: "POST",
        body: { email, password },
      }),
    }),

    register: builder.mutation<RegisterResponse, RegisterRequest>({
      query: ({ nickname, email, password }) => ({
        url: "/register",
        method: "POST",
        body: { nickname, email, password },
      }),
    }),
  }),
});

export const { useLoginMutation, useRegisterMutation } = authApi;
