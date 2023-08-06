import React from "react";
import { Router } from "./lib/router/Router";
import { Store } from "./lib/feature/Store";
import { Theme } from "./lib/theme/Theme";
import "./App.css";

export const App: React.FC = () => {
  return (
    <Theme>
      <Store>
        <Router />
      </Store>
    </Theme>
  );
};
