import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { IChannel } from "../entity/entities";

interface GetListChannelsRequest {
  token: string;
}

interface GetListChannelsResponse {
  channels: IChannel[];
}

interface CreateChannelRequest {
  token: string;
  title: string;
}

interface CreateChannelResponse {
  channel: IChannel;
}

interface DeleteChannelRequest {
  token: string;
  id: string;
}

interface DeleteChannelResponse {}

export const channelsApi = createApi({
  reducerPath: "channelsApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/channels`,
  }),
  endpoints: (builder) => ({
    getListChannels: builder.query<
      GetListChannelsResponse,
      GetListChannelsRequest
    >({
      query: ({ token }) => ({
        url: "/",
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),

    createChannel: builder.mutation<
      CreateChannelResponse,
      CreateChannelRequest
    >({
      query: ({ token, title }) => ({
        url: "/",
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: {
          title,
        },
      }),
    }),

    deleteChannel: builder.mutation<
      DeleteChannelResponse,
      DeleteChannelRequest
    >({
      query: ({ token, id }) => ({
        url: `/${id}`,
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }),
    }),
  }),
});

export const {
  useGetListChannelsQuery,
  useCreateChannelMutation,
  useDeleteChannelMutation,
} = channelsApi;
