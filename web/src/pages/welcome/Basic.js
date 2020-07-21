import React from 'react';
import { Button, Form, Input, message } from 'antd';
import styles from './index.module.css';
import {httpPatch} from '@/utils/request';

class Basic extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      nickname: '',
    }
  }

  componentDidMount() {
    this.setState({
      nickname: JSON.parse(sessionStorage.getItem('user')).nickname
    })
  };

  handleSubmit = () => {
    if (this.state.nickname == "") {
      message.warn("昵称不能为空")
      return
    }
    if (JSON.parse(sessionStorage.getItem('user')).nickname == this.state.nickname) {
      message.warn("昵称没有改变")
      return
    }
    this.setState({loading: true})
    httpPatch('/admin/user', {nickname: this.state.nickname, type: "nickname"})
      .then(() => {
        message.success('设置成功，重新登录或刷新页面后生效');
        var data = JSON.parse(sessionStorage.getItem('user'))
        data.nickname = this.state.nickname
        sessionStorage.setItem('user', JSON.stringify(data))
      })
      .finally(() => this.setState({loading: false}))
  }

  render() {
    return (
      <React.Fragment>
        <div className={styles.title}>基本设置</div>
        <Form style={{maxWidth: 320}}>
          <Form.Item colon={false} label="昵称">
            <Input value={this.state.nickname} placeholder="请输入" onChange={e => this.setState({nickname: e.target.value})}/>
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={this.state.loading} onClick={this.handleSubmit}>保存设置</Button>
          </Form.Item>
        </Form>
      </React.Fragment>
    )
  }
}

export default Basic