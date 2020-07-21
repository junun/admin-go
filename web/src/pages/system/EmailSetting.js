import React from 'react';
import {Button, Form, Input, Radio, message, Popover} from 'antd';
import {connect} from "dva";
import styles from './index.module.css';
import lds from 'lodash';

class EmailSetting extends React.Component {
  constructor(props) {
    super(props);
    // this.setting = JSON.parse(lds.get(store.settings, 'mail_service.value', "{}"));
    this.setting = JSON.parse(lds.get(this.props.settings, 'mail_service.Value', "{}"));
    this.state = {
      loading: false,
      mode: this.setting['server'] === undefined ? '1' : '2',
      spug_key: lds.get(this.props.settings, 'spug_key.value', ""),
      mail_test_loading: false,
    }
  }

  handleEmailTest = () => {
    this.props.form.validateFields((error, data) => {
      console.log
      if (!error) {
        this.setState({mail_test_loading: true});
        this.props.dispatch({ 
          type: 'user/settingMailTest',
          payload: data,
        }).finally(()=> this.setState({mail_test_loading: false}))
      }
    })
  };

  handleSubmit = () => {
    const formData = [];
    this.props.form.validateFields((err, data) => {
      if (!err) {
        this.setState({loading: true});
        formData.push({name: 'mail_service', value: JSON.stringify(data)});
        this.props.dispatch({ 
          type: 'user/settingModify',
          payload: {
            Data: formData,
          }
        }).finally(() => this.setState({loading: false}))
      }
    })
  };

  render() {
    const {getFieldDecorator} = this.props.form;
    const {loading} = this.state;
    return (
      <React.Fragment>
        <div className={styles.title}>邮件服务设置</div>
        <Form style={{maxWidth: 340}}>
          <Form.Item colon={false} label="邮件服务" help="用于通过邮件方式发送报警信息">
              <Form.Item labelCol={{span: 8}} wrapperCol={{span: 16}} required label="邮件服务器">
                {getFieldDecorator('server', {
                  initialValue: this.setting['server'], rules: [
                    {required: true, message: '请输入邮件服务器地址'}
                  ]
                })(
                  <Input placeholder="例如：smtp.exmail.qq.com"/>
                )}
              </Form.Item>
              <Form.Item labelCol={{span: 8}} wrapperCol={{span: 16}} required label="端口">
                {getFieldDecorator('port', {
                  initialValue: this.setting['port'], rules: [
                    {required: true, message: '请输入邮件服务端口'}
                  ]
                })(
                  <Input placeholder="例如：465"/>
                )}
              </Form.Item>
              <Form.Item labelCol={{span: 8}} wrapperCol={{span: 16}} required label="邮箱账号">
                {getFieldDecorator('username', {
                  initialValue: this.setting['username'], rules: [
                    {required: true, message: '请输入邮箱账号'}
                  ]
                })(
                  <Input placeholder="例如：dev@exmail.com"/>
                )}
              </Form.Item>
              <Form.Item labelCol={{span: 8}} wrapperCol={{span: 16}} required label="密码/授权码">
                {getFieldDecorator('password', {
                  initialValue: this.setting['password'], rules: [
                    {required: true, message: '请输入邮箱账号对应的密码或授权码'}
                  ]
                })(
                  <Input.Password placeholder="请输入对应的密码或授权码"/>
                )}
              </Form.Item>
              <Form.Item labelCol={{span: 8}} wrapperCol={{span: 16}} label="发件人昵称">
                {getFieldDecorator('nickname', {initialValue: this.setting['nickname']})(
                  <Input placeholder="请输入发件人昵称"/>
                )}
              </Form.Item>
          </Form.Item>
          <div>
          <Button
            type="danger" loading={this.state.mail_test_loading}  style={{ marginRight: '10px' }}
            onClick={this.handleEmailTest}>测试邮件服务</Button>
          <Button
            type="primary" loading={loading} style={{ marginTop: 20}}
            onClick={this.handleSubmit}>保存设置</Button>
          </div>
        </Form>
      </React.Fragment>
    )
  }
}

export default Form.create()(EmailSetting)
