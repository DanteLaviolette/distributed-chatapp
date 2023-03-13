import { useState } from "react";
import { getSignedInUser } from "./utils";

function App() {
  const [user, setUser] = useState(getSignedInUser())
  return (
    <div>
      <p>Hello World!</p>
    </div>
  );
}

export default App;
