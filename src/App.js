import './App.css';

import { Amplify, API } from 'aws-amplify';
import awsExports from './aws-exports';
import { Authenticator } from '@aws-amplify/ui-react';
import '@aws-amplify/ui-react/styles.css';
import { useEffect, useState } from 'react';

Amplify.configure(awsExports);

function App() {
  return (
    <Authenticator variation='modal'>
      {({ signOut, user }) => (
        <div className="App">
          <h1>Ranking</h1>
          <RankingList />
        </div>
      )}
    </Authenticator>
  );
}

function RankingList() {
  const [players, setPlayers] = useState([])

  useEffect(() => {
    const apiName = 'api386f423e';
    const path = '/getRankings';

    API
      .get(apiName, path)
      .then(response => {
        setPlayers(response)
      })
      .catch(error => {
        console.log(error.response);
      });
  }, [])

  if (players.length !== 0) {
    var players_sorted = [...players]
    players_sorted.sort((a, b) => {
      if (a["Grade"][0] === b["Grade"][0]) {
        return a["Grade"].charCodeAt(1) - b["Grade"].charCodeAt(1)
      } else {
        return a["Grade"].charCodeAt(0) - b["Grade"].charCodeAt(0)
      }
    })

    return (
      players_sorted.map((p, i) => {
        return (<div key={i}>{p["Name"] + "  " + p["Grade"]}</div>)
      })
    )
  } else {
    return (
      <div>
        Loading
      </div>
    )
  }
}

export default App;
