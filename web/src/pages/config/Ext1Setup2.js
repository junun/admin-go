import React from 'react';
import { Form, Input, InputNumber, Select, Button, Icon, Col } from "antd";
import { Link } from 'react-router-dom';
import styles from './index.module.css';

class Ext1Setup2 extends React.Component {
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

  showHostName = (value) => {
    var tmp = [];
    var ids = this.props.info["HostIds"].split(",");
    this.props.hostListByAppId.map ((item) => {
      if (ids.includes(item.id.toString())) {
        tmp.push(item.Name)
      }
    });

    return tmp.join(",")
  }

  onInputDstDirChange = (e) => {
    var tmp = this.state.info;
    tmp['DstDir'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  onInputDstRepoChange = (e) => {
    var tmp = this.state.info;
    tmp['DstRepo'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  onInputVersionsChange = (e) => {
    var tmp = this.state.info;
    tmp['Versions'] = e;
    this.setState({ 
      info: tmp,
    });
  };
  
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
    return info['DstDir'] && info['DstRepo'] && info['Versions'] && info['HostIds']!=""
  };

  render() {
    const {info, defaultValue} = this.state;
    const hostList = this.props.hostListByAppId;

    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="目标主机部署路径" help="目标主机的应用根目录，例如：/var/www/html">
          <Input value={info['DstDir']} onChange={this.onInputDstDirChange} placeholder="请输入目标主机部署路径"/>
        </Form.Item>
        <Form.Item required label="目标主机仓库路径" help="此目录用于存储应用的历史版本，例如：/data/spug/repos">
          <Input value={info['DstRepo']} onChange={this.onInputDstRepoChange} placeholder="请输入目标主机仓库路径"/>
        </Form.Item>
        <Form.Item required label="保留历史版本数量" help="早于指定数量的历史版本会被删除，以释放空间">
          <InputNumber value={info['Versions']} onChange={this.onInputVersionsChange} placeholder="10"/>
        </Form.Item>
        <Form.Item required label="发布目标主机">
          <Col span={16}>
            <Select
              mode="multiple"
              style={{ width: '100%' }}
              placeholder={this.props.info['HostIds']&&this.showHostName(this.props.info['HostIds'])||"请选择"}
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

export default Ext1Setup2