import { useState } from "react";
import TopBar from "./TopBar/TopBar";
import { getSignedInUser } from "./utils/utils";
import Chat from "./Chat/Chat"

function App() {
  // Set signed in user in state -- this can be used/updated by other components
  const [user, setUser] = useState(getSignedInUser())
  return (
    <div style={{width: "100%", height: "100vh", display: 'flex', flexDirection: 'column'}}>
      <TopBar user={user} setUser={setUser}/>
      <Chat user={user} setUser={setUser} />
    </div>
  );
}

export default App;
