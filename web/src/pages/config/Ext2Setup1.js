import React from 'react';
import { Link } from 'react-router-dom';
import { Switch, Col, Form, Select, Button, Input } from "antd";

const tmpObj = {id: 0, Name: "关闭"}
class Ext2Setup1 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      info: {},
      localRobot: [],
    }
  }

  componentDidMount() {
    console.log(this.props.robotList);
    // 兼容 关闭的情况
    var tmpArr = this.props.robotList;
    tmpArr.push(tmpObj);
    this.setState({ 
      info: this.props.info,
      localRobot: tmpArr,
    });
  }

  onSwitchChange = (e) => {
    var tmp = this.state.info;
    tmp['EnableCheck'] = e;
    this.setState({ 
      info: tmp,
    });
  };


  onTnameInputChange = (e) => {
    var tmp = this.state.info;
    tmp['TemplateName'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };


  handleNotifyChange = (e) => {
    var tmp = this.state.info;
    tmp['NotifyId'] = parseInt(e);
    this.setState({ 
      info: tmp,
    });
  };

  nextPage = () => {
    this.props.nextPage(this.state.info)
  };


  render() {
    const {info, localRobot} = this.state;
    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="发布模板名字">
          <Input value={info['TemplateName']} onChange={this.onTnameInputChange} placeholder="请输入发布模板名字"/>
        </Form.Item>
        <Form.Item label="发布审核">
          <Switch
            checkedChildren="开启"
            unCheckedChildren="关闭"
            checked={info['EnableCheck'] && true || false}
            onChange={this.onSwitchChange}/>
        </Form.Item>
        <Form.Item label="结果通知" help="应用发布成功或失败结果通知">
          <Col span={16}>
            <Select
              value={info['NotifyId'] || 0}
              style={{ width: '100%' }}
              onChange={this.handleNotifyChange}
            >
              {localRobot.map(x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
            </Select>
          </Col>
          <Col span={6} offset={2}>
            <Link to="/system/robot">新建机器人通道</Link>
          </Col>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button
            type="primary"
            disabled={!(info['TemplateName'])}
            onClick={this.nextPage}>下一步</Button>
        </Form.Item>
      </Form>
    )
  }
}

export default Ext2Setup1
