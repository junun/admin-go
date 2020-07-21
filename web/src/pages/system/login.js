import React from 'react';
import { connect } from 'dva';
import {
  Form, Icon, Input, Button, Checkbox
} from 'antd';

import styles from './login.css';

@connect(({ loading }) => {
  return {
    loading: loading.global,
  }
})
class NormalLoginForm extends React.Component {
  handleSubmit = (e) => {
    e.preventDefault();
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        dispatch({
          type: 'user/login',
          payload: values,
        });
      }
    });
  }

  render() {
    const { getFieldDecorator } = this.props.form;
    // const {loginType} = this.props.state;
    // console.log(loginType)
    return (
      <Form onSubmit={this.handleSubmit} className={styles.login_form}>
        <Form.Item>
          {getFieldDecorator('username', {
            rules: [{ required: true, message: '请输入用户名' }],
          })(
            <Input prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />} placeholder="用户名" />
          )}
        </Form.Item>
        <Form.Item>
          {getFieldDecorator('password', {
            rules: [{ required: true, message: '请输入密码' }],
          })(
            <Input prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />} type="password" placeholder="密码" />
          )}
        </Form.Item>
        <Form.Item>
          {getFieldDecorator('remember', {
            valuePropName: 'checked',
            initialValue: true,
          })(
            <Checkbox>记住用户名</Checkbox>
          )}
          {/* <a className={styles.login_form_forgot} href="">Forgot password</a> */}
          
          <Button loading={this.props.loading} type="primary" htmlType="submit" className={styles.login_form_button}>
            登录
          </Button>
          {/* Or <a href="">register now!</a> */}
        </Form.Item>
      </Form>
    );
  }
}

export default Form.create()(NormalLoginForm);
