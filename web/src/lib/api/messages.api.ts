import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { IMessage } from "../entity/entities";

interface GetListMessagesRequest {
  channelId: string;
  token: string;
}

interface GetListMessagesResponse {
  messages: IMessage[];
}

export const messagesApi = createApi({
  reducerPath: "messagesApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/channels`,
  }),
  endpoints: (builder) => ({
    getListMessages: builder.query<
      GetListMessagesResponse,
      GetListMessagesRequest
    >({
      query: ({ channelId, token }) => ({
        url: `/${channelId}/messages`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),
  }),
});

export const { useGetListMessagesQuery } = messagesApi;
