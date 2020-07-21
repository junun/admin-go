import React from 'react';
import { Alert, Button, Form, Input, Modal, message } from 'antd';
import styles from './index.module.css';
import lds from 'lodash';

class KeySetting extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      settings: {},
    }
  }

  componentDidMount() {
    this.setState({
      settings: this.props.settings
    })
  }

  handlePublicKey = (e) => {
    var tmp = this.props.settings
    lds.set(tmp, 'public_key.Value', e.target.value)
    this.setState({
      settings: tmp
    })
  }

  handlePrivateKey = (e) => {
    var tmp = this.props.settings
    lds.set(tmp, 'private_key.Value', e.target.value)
    this.setState({
      settings: tmp
    })
  }

  handleSubmit = () => {
    Modal.confirm({
      title: '密钥修改确认',
      content: <span style={{color: '#f5222d'}}>请谨慎修改密钥对，修改密钥对会让现有的主机都无法进行验证，影响与主机相关的各项功能！</span>,
      onOk: () => {
        Modal.confirm({
          title: '小提示',
          content: <div>修改密钥对需要<span style={{color: '#f5222d'}}>重启服务后生效</span>，已添加的主机需要重新进行编辑验证后才可以正常连接。</div>,
          onOk: this.doModify
        })
      }
    })
  }

  doModify = () => {
    const public_key  = lds.get(this.state.settings, 'public_key.Value');
    const private_key = lds.get(this.state.settings, 'private_key.Value');
    this.props.dispatch({ 
      type: 'user/settingModify',
      payload: {
        Data: [{Name: 'public_key', Value: public_key}, {Name: 'private_key', Value: private_key}],
      }
    }).finally(() => this.setState({loading: false}))
  }

  render() {
    const {loading, settings} = this.state;
    return (
      <React.Fragment>
        <div className={styles.title}>密钥设置</div>
        <Alert
          closable
          showIcon
          type="info"
          style={{width: 650}}
          message="小提示"
          description="在你没有上传密钥的情况下，Spug会在首次添加主机时自动生成密钥对。"
        />
        <Form style={{maxWidth: 650}}>
          <Form.Item label="公钥" help="一般位于 ~/.ssh/id_rsa.pub">
            <Input.TextArea
              rows={7}
              spellCheck={false}
              value={lds.get(settings, 'public_key.Value', "")}
              onChange={e => this.handlePublicKey(e)}
              placeholder="请输入公钥"/>
          </Form.Item>
          <Form.Item label="私钥" help="一般位于 ~/.ssh/id_rsa">
            <Input.TextArea
              rows={14}
              spellCheck={false}
              value={lds.get(settings, 'private_key.Value', "")}
              onChange={e => this.handlePrivateKey(e)}
              placeholder="请输入私钥"/>
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={loading} onClick={this.handleSubmit}>保存设置</Button>
          </Form.Item>
        </Form>
      </React.Fragment>
    )
  }
}

export default KeySetting
