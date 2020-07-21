import React from 'react';
import {connect} from "dva";
import { Form, Input, InputNumber, Select, Button, Icon, Col } from "antd";
import { Link } from 'react-router-dom';
import styles from './index.module.css';

@connect(({ loading, host }) => {
  return {
    hostListByAppId: host.hostListByAppId,
  }
})

class Ext2Setup2 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      info: {},
    }
  }

  componentDidMount() {
    this.props.dispatch({ 
      type: 'host/getHostByAppId',
      payload: {
        aid: this.props.info['Aid'],
      }
    })

    this.setState({ 
      info: this.props.info,
    });
  }

  handleHostChange = (e) => {
    var tmp = this.state.info;
    tmp['HostIds'] = "";
    e.map ((item, index) =>  {
      if (index > 0) {
        tmp['HostIds'] = tmp['HostIds'] + "," + item
      } else {
        tmp['HostIds'] += item
      }
    });

    this.setState({ 
      info: tmp,
    });
  };

  nextPage = () => {
    this.props.nextPage(this.state.info)
  };

  prePage = () => {
    this.props.prePage(this.state.info)
  };

  checkStatus = () => {
    const info = this.state.info;
    return info['HostIds'] != undefined && info['HostIds'] != ""
  };

  render() {
    const {info, defaultValue} = this.state;
    const hostList = this.props.hostListByAppId;

    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="发布目标主机">
          <Col span={16}>
            <Select
              mode="multiple"
              style={{ width: '100%' }}
              value={this.props.info['HostIds'] && JSON.parse("[" + this.props.info['HostIds'] + "]") || []}
              filterOption={(input, option) => option.props.children[0].toLowerCase().indexOf(input.toLowerCase()) >= 0}
              onChange={this.handleHostChange}
            >
              {hostList.map(item => (
                <Select.Option key={item.id} value={item.id}>
                  {item.Name}({item['Addres']}:{item['Port']})
                </Select.Option>
              ))}
            </Select>
          </Col>
          <Col span={6} offset={2}>
            <Link to="/host/list">业务绑定主机</Link>
          </Col>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button disabled={!this.checkStatus()} type="primary" onClick={this.nextPage}>下一步</Button>
          <Button style={{marginLeft: 20}} onClick={this.prePage}>上一步</Button>
        </Form.Item>
      </Form>
    )
  }
}

export default Ext2Setup2