import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/dist/query";
import { IChannel } from "../entity/entities";

interface GetListChannelsRequest {
  token: string;
}

interface GetListChannelsResponse {
  channels: IChannel[];
}

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
  }),
});
