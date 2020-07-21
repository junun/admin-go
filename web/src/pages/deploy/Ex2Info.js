import React, {Fragment, Component} from 'react';
import {connect} from "dva"
import { Link } from 'react-router-dom';
import { Modal, Form, Select, Button, Icon, Input, Col } from 'antd';
import styles from './index.module.css';

@connect(({ loading, deploy, config }) => {
  return {
  }
})

class Ex2Info extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      info: {},
    }
  }

  componentDidMount() {
    this.setState({
      info: this.props.info,
    })
  }

  handleSubmit = () => {
    this.setState({loading: true});
    console.log(this.state.info);
    this.props.dispatch({
      type: 'deploy/deployAdd',
      payload: this.state.info,
    }).then(() => {
      this.setState({loading: false});
      this.props.onCancel();
    })
  };

  prePage = () => {
    this.props.prePage(this.state.info)
  };

  onDescInputChange = (e) => {
    var tmp = this.props.info;
    tmp['Desc'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };


  render() {
    const {info} = this.state;

    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item label="备注信息">
            <Input value={ info['Desc']} onChange={this.onDescInputChange} placeholder="请输入备注信息"/>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button type="primary" onClick={this.handleSubmit}>提交</Button>
          <Button style={{marginLeft: 20}} onClick={this.prePage}>上一步</Button>
        </Form.Item>
      </Form> 
    )
  }
}

export default Ex2Info
