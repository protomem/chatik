import { Box, Link, Typography } from "@mui/joy";
import React from "react";

export const NotFoundPage: React.FC = () => {
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        mt: 10,
      }}
    >
      <Typography component={"h1"} fontSize={"3rem"}>
        Page Not Found
      </Typography>
      <Link fontSize={"1.5rem"} href={"/"} variant={"solid"} underline={"none"}>
        Go Home
      </Link>
    </Box>
  );
};
