import React from 'react';
import { Modal, List, Icon } from 'antd';
import styles from './index.module.css';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.socket   = null;
    this.state = {
      data: [],
    };
  }

  state = {
    loading: false,
    hasMore: true,
  };

  handleInfiniteOnLoad = () => {
    let { data } = this.state;
    this.setState({
      loading: true,
    });
    if (data.length > 14) {
      message.warning('Infinite List loaded all');
      this.setState({
        hasMore: false,
        loading: false,
      });
      return;
    }
  };

  componentDidMount() {
    const token = sessionStorage.getItem('jwt');
    const id = this.props.id ? this.props.id: 0;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    
    this.socket = new window.WebSocket(`${protocol}//127.0.0.1:8080/admin/undeploy/ws/${id}/ssh/${token}`);

    var thus = this;
    thus.socket.onopen = function () {
      var tmp = thus.state.data
      tmp.push("websock onopen！");
      thus.setState({
        data : tmp,
      })
    };

    var thus = this;
    thus.socket.onclose = function (evt) {
      var tmp = thus.state.data
      tmp.push("End");
      thus.setState({
        data : tmp,
      })
    };

    // var thus = this;
    thus.socket.onmessage = e => {
      if (e.data === 'pong') {
        thus.socket.send(JSON.stringify({type: "heartbeat", data: ""}));
      } else {
        var tmp = thus.state.data
        tmp.push(e.data);
        thus.setState({
          data : tmp,
        })
      }
    };

    if (this.refs.chatoutput != null) {
      this.refs.chatoutput.scrollTop = this.refs.chatoutput.scrollHeight;
    }
  }

  componentWillUnmount() {
    this.socket.close()
  }

  componentDidUpdate() {
    if (this.refs.chatoutput != null) {
      this.refs.chatoutput.scrollTop = this.refs.chatoutput.scrollHeight;
    }
  }

  render() {
    return (
      <Modal
        visible
        width={800}
        destroyOnClose= "true"
        title="回滚信息控制台"
        footer={null}
        onCancel={this.props.onCancel}
        // onOk={this.handleSubmit}
        maskClosable={false}
        className={styles.modal}>
        <div ref='chatoutput' className={styles.modaldiv}>
          {this.state.data.map((item, index) => <div key={index}>{item}</div>)}
        </div>
      </Modal>
    )
  }
}