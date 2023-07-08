import { Router } from "./lib/router/Router";
import { Theme } from "./lib/theme/Theme";

export function App() {
  return (
    <Theme>
      <Router />
    </Theme>
  );
}
