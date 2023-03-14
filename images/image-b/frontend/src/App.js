import { useState } from "react";
import TopBar from "./TopBar/TopBar";
import { getSignedInUser } from "./utils";

function App() {
  const [user, setUser] = useState(getSignedInUser())
  return (
    <div>
      <TopBar user={user} setUser={setUser}/>
      <p>Hello World!</p>
    </div>
  );
}

export default App;
