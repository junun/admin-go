import React from 'react';
import { Link } from 'react-router-dom';
import { Switch, Col, Form, Input, Select, Button } from "antd";

class Ext1Setup1 extends React.Component {
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

  onSwitchChange = (e) => {
    var tmp = this.props.info;
    tmp['EnableCheck'] = e;
    this.setState({ 
      info: tmp,
    });
  };

  onInputChange = (e) => {
    var tmp = this.props.info;
    tmp['RepoUrl'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  onTnameInputChange = (e) => {
    var tmp = this.props.info;
    tmp['TemplateName'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  onTagInputChange = (e) => {
    var tmp = this.props.info;
    tmp['Tag'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  nextPage = () => {
    this.props.nextPage(this.state.info)
  };

  render() {
    const {info } = this.state;
    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="发布模板名字">
          <Input value={info['TemplateName']} onChange={this.onTnameInputChange} placeholder="请输入发布模板名字"/>
        </Form.Item>
        <Form.Item required label="Git仓库地址">
          <Input value={info['RepoUrl']} onChange={this.onInputChange} placeholder="请输入Git仓库地址"/>
        </Form.Item>
        <Form.Item label="编译版本号">
          <Input value={info['Tag']} onChange={this.onTagInputChange} placeholder="请输入Tag版本号，如1.0.0"/>
        </Form.Item>
        <Form.Item label="发布审核">
          <Switch
            checkedChildren="开启"
            unCheckedChildren="关闭"
            checked={info['EnableCheck'] && true || false}
            onChange={this.onSwitchChange}/>
        </Form.Item>

        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button
            type="primary"
            disabled={!(info['TemplateName'] && info['RepoUrl'])}
            onClick={this.nextPage}>下一步</Button>
        </Form.Item>
      </Form>
    )
  }
}

export default Ext1Setup1