import React, {Fragment, Component} from 'react';
import { Link } from 'react-router-dom';
import { Modal, Form, Select, Button, Icon, Input, Col, Steps } from 'antd';
import styles from './index.module.css';

class Template extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      info: {},
    }
  }

  componentDidMount() {
    this.setState({ 
      info: this.props.info,
    });
  }

  onNameInputChange = (e) => {
    var tmp = this.state.info;
    tmp['Name'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  handleTidChange = (e) => {
    var tmp = this.state.info;
    if (tmp['Tid'] != parseInt(e)) {
      tmp['originTid'] = tmp['Tid']
      tmp['Tid'] = parseInt(e)
    }

    this.setState({ 
      info: tmp,
    });
  };

  handler = (values) => {
    var tmpExtend = 1
    for (var i = 0; i < this.props.appTemplateList.length; i++) {
      if (this.props.appTemplateList[i].Dtid === values.Tid) {
        tmpExtend = this.props.appTemplateList[i].Extend
        break
      }
    }

    values['Extend'] = tmpExtend
    
    this.props.nextPage(values)
  }

  // nextPage = () => {
  //   this.props.nextPage(this.state.info)
  // };

  render() {
    const {info} = this.state;
    const {appTemplateList} = this.props;

    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="上线单标题">
          <Input value={ info['Name']} onChange={this.onNameInputChange} placeholder="请输入上线单名字"/>
        </Form.Item>
        <Form.Item required label="择项目发布模板">
          <Col span={16}>
            <Select
              value={ info['Tid']}
              style={{ width: '100%' }}
              onChange={this.handleTidChange}
            >
              {appTemplateList.map(x => <Select.Option key={x.Dtid} value={x.Dtid}>{x.TemplateName}</Select.Option>)}
            </Select>
          </Col>
          <Col span={6} offset={2}>
            <Link to="/config/app">新建发布模板</Link>
          </Col>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button
            type="primary"
            disabled={!( info['Name'] && info['Tid'])}
            onClick={() => this.handler(this.state.info)}>下一步</Button>
        </Form.Item>
      </Form> 
    )
  }
}

export default Template
