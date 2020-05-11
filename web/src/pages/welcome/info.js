import React, { Component } from 'react';
import { Card } from 'antd';

class WelcomePage extends Component {
  componentDidMount() {
  }

  render() {
    return (
      <Card>
        <div>{JSON.parse(sessionStorage.getItem('user')).nickname}, 欢迎你</div>

      </Card>
    )
  }
}

export default WelcomePage;