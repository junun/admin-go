import React, { useState } from 'react';
import { Button, Form, Input, message } from 'antd';
import styles from './index.module.css';
import {httpPatch, httpPost} from '@/utils/request';
import router from 'umi/router';
import history from '@/utils/history';

class Reset extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      old_password: "",
      new_password: "",
      new2_password: "",
    }
  }

  handleSubmit = () => {
    if (this.state.old_password == "") {
      return message.error('请输入原密码')
    }
    if (this.state.new_password == "") {
      return message.error('请输入新密码')
    }
    if (this.state.new_password != this.state.new2_password) {
      return message.error('两次输入密码不一致')
    }
    if (this.state.new_password.length < 6) {
      return message.error('请设置至少6位的新密码')
    }

    this.setState({loading: true})
    httpPatch('/admin/user', {
      old_password: this.state.old_password,
      new_password: this.state.new_password,
      type: "password"}).then(res => {
        if (res.code == 200) {
          message.success(res.message);
          httpPost('/admin/user/logout')
          router.push('/user/login')
        } else {
          message.error(res.message);
        }
      }).finally(() => this.setState({loading: false}))
  }

  render() {
    return (
      <React.Fragment>
        <div className={styles.title}>修改密码</div>
        <Form style={{maxWidth: 320}} labelCol={{span: 6}} wrapperCol={{span: 18}}>
          <Form.Item label="原密码">
            <Input.Password value={this.state.old_password} placeholder="请输入" onChange={e => this.setState({old_password:e.target.value})}/>
          </Form.Item>
          <Form.Item label="新密码">
            <Input.Password value={this.state.new_password} placeholder="请输入" onChange={e => this.setState({new_password:e.target.value})}/>
          </Form.Item>
          <Form.Item label="再次确认">
            <Input.Password value={this.state.new2_password} placeholder="请输入" onChange={e => this.setState({new2_password:e.target.value})}/>
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={this.state.loading} onClick={this.handleSubmit}>保存设置</Button>
          </Form.Item>
        </Form>
      </React.Fragment>
    )
  }
}

export default Reset
