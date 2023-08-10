import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { MessageEntity } from "../../entities/entities";

interface GetAllMessagesRequest {
  channelId: string;
  accessToken: string;
}

interface GetAllMessagesResponse {
  messages: MessageEntity[];
}

export const messagesApi = createApi({
  reducerPath: "messagesApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/channels`,
  }),
  endpoints: (builder) => ({
    getAllMessages: builder.query<
      GetAllMessagesResponse,
      GetAllMessagesRequest
    >({
      query: ({ channelId, accessToken }) => ({
        url: `/${channelId}/messages`,
        method: "GET",
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }),
    }),
  }),
});

export const { useGetAllMessagesQuery } = messagesApi;
