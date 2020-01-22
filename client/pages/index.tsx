import React from 'react';
import axios from 'axios';
import fetch from 'isomorphic-unfetch'

import { NextPage } from 'next';

const Home: NextPage = (props: any) => {
  const [token, setToken] = React.useState('');
  React.useEffect(() => {
    axios
      .post('api/sessions', {
        email: 'user2@example.org',
        password: 'password',
      })
      .then(res => setToken(res.data));
  }, []);
  return (
    <>
      <h1>Get users works!! </h1>
      <pre>{JSON.stringify(props.data, null, ' ')}</pre>
      <span style={{ flexWrap: 'wrap' }}>
        Token: {JSON.stringify(token.split('.'), null, '\t')}
      </span>
    </>
  );
};

Home.getInitialProps = async ({ req }) => {
  const res = await fetch('http://nginx/api/sessions', {
    method: 'POST',
    body: JSON.stringify({
      email: 'user2@example.org',
      password: 'password',
    }),
  });
  const data = await res.json();
  return await { data };
};

export default Home;
