import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { IMessage } from "../entity/entities";

interface GetListMessagesRequest {
  channelId: string;
  token: string;
}

interface GetListMessagesResponse {
  messages: IMessage[];
}

interface CreateMessageRequest {
  channelId: string;
  token: string;
  content: string;
}

interface CreateMessageResponse {
  message: IMessage;
}

interface DeleteMessageRequest {
  channelId: string;
  messageId: string;
  token: string;
}

interface DeleteMessageResponse {}

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

    createMessage: builder.mutation<
      CreateMessageResponse,
      CreateMessageRequest
    >({
      query: ({ channelId, token, content }) => ({
        url: `/${channelId}/messages`,
        method: "POST",
        body: { content },
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),

    deleteMessage: builder.mutation<
      DeleteMessageResponse,
      DeleteMessageRequest
    >({
      query: ({ channelId, messageId, token }) => ({
        url: `/${channelId}/messages/${messageId}`,
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),
  }),
});

export const {
  useGetListMessagesQuery,
  useCreateMessageMutation,
  useDeleteMessageMutation,
} = messagesApi;
