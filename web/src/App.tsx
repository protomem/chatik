import { Router } from "./lib/router/Router";
import { StoreProvider } from "./lib/store/StoreProvider";
import { Theme } from "./lib/theme/Theme";

export function App() {
  return (
    <Theme>
      <StoreProvider>
        <Router />
      </StoreProvider>
    </Theme>
  );
}
