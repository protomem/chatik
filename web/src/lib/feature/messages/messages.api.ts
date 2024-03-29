import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { MessageEntity } from "../../entities/entities";

interface GetAllMessagesRequest {
  channelId: string;
  accessToken: string;
}

interface GetAllMessagesResponse {
  messages: MessageEntity[];
}

interface CreateMessageRequest {
  channelId: string;
  content: string;
  accessToken: string;
}

interface CreateMessageResponse {
  message: MessageEntity;
}

interface DeleteMessageRequest {
  channelId: string;
  messageId: string;
  accessToken: string;
}

interface DeleteMessageResponse {}

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

    createMessage: builder.mutation<
      CreateMessageResponse,
      CreateMessageRequest
    >({
      query: ({ channelId, content, accessToken }) => ({
        url: `/${channelId}/messages`,
        method: "POST",
        body: { content },
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }),
    }),

    deleteMessage: builder.mutation<
      DeleteMessageResponse,
      DeleteMessageRequest
    >({
      query: ({ channelId, messageId, accessToken }) => ({
        url: `/${channelId}/messages/${messageId}`,
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }),
    }),
  }),
});

export const {
  useGetAllMessagesQuery,
  useCreateMessageMutation,
  useDeleteMessageMutation,
} = messagesApi;
