import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { ChannelEntity } from "../../entities/entities";

interface GetAllChannelsRequest {
  accessToken: string;
}

interface GetAllChannelsResponse {
  channels: ChannelEntity[];
}

interface CreateChannelRequest {
  title: string;
  accessToken: string;
}

interface CreateChannelResponse {
  channel: ChannelEntity;
}

export const channelsApi = createApi({
  reducerPath: "channelsApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api/v1/channels`,
  }),
  endpoints: (builder) => ({
    getAllChannels: builder.query<
      GetAllChannelsResponse,
      GetAllChannelsRequest
    >({
      query: ({ accessToken }) => ({
        url: "/",
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }),
    }),

    createChannel: builder.mutation<
      CreateChannelResponse,
      CreateChannelRequest
    >({
      query: ({ title, accessToken }) => ({
        url: "/",
        method: "POST",
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        body: {
          title,
        },
      }),
    }),
  }),
});

export const { useGetAllChannelsQuery, useCreateChannelMutation } = channelsApi;
