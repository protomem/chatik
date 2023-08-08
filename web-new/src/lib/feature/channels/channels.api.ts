import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { ChannelEntity } from "../../entities/entities";

interface GetAllChannelsRequest {
  accessToken: string;
}

interface GetAllChannelsResponse {
  channels: ChannelEntity[];
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
  }),
});

export const { useGetAllChannelsQuery } = channelsApi;
