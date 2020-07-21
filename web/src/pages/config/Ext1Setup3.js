import React from 'react';

import { Form, Row, Col, Button, Radio, Icon, message } from "antd";
import Editor from 'react-ace';
import 'ace-builds/src-noconflict/mode-text';
import 'ace-builds/src-noconflict/mode-sh';
import 'ace-builds/src-noconflict/theme-tomorrow';
import styles from './index.module.css';
import {cleanCommand} from "@/utils/globalTools";

class Ext1Setup3 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      info: {},
      full: '',
      filterRule: {},
    }
  }

  componentDidMount() {
    var tmpFilterRule = {}
    if (this.props.info["FilterRule"] != undefined && this.props.info["FilterRule"] != "") {
      tmpFilterRule = JSON.parse(this.props.info["FilterRule"])
    }

    var tmpCustomEnvs = {}
    if (this.props.info["CustomEnvs"] != undefined && this.props.info["CustomEnvs"] != "") {
      tmpCustomEnvs = JSON.parse(this.props.info["CustomEnvs"])
    }

    var strCustomEnvs = ""
    for(let key in tmpCustomEnvs){
        strCustomEnvs = strCustomEnvs + key + "=" + tmpCustomEnvs[key] + "\n"
    } 

    var tmp = this.props.info
    tmp.CustomEnvs = strCustomEnvs
    this.setState({ 
      info: tmp,
      filterRule: tmpFilterRule,
    });
  }

  handleSubmit = () => {
    var values = this.state.info
    values.EnableCheck  = values.EnableCheck ? 1 : 0
    values.FilterRule   = JSON.stringify(this.state.filterRule)
    values.Extend = 1
    if (values.Dtid === undefined) {
      this.props.dispatch({
        type: 'config/deployExtendAdd',
        payload: values,
      }).then(()=> {
        this.props.handleCancel()
      });
    } else {
      this.props.dispatch({
        type: 'config/deployExtendEdit',
        payload: values,
      }).then(()=> {
        this.props.handleCancel()
      });
    }
  };

  prePage = () => {
    this.props.prePage(this.state.info)
  };

  setStateFull = (id)  => {
    this.setState({ 
      full: id,
    });
  }

  onChangeFilterRuleType = (e) => {
    var tmp = this.state.info;
    var tmpfilterRule = this.state.filterRule
    tmpfilterRule["type"] = e.target.value
    this.setState({ 
      info: tmp,
      filterRule: tmpfilterRule,
    })
  };

  FilterLabel = (props) => (
    <div style={{display: 'inline-block', height: 39, width: 344}}>
      <span style={{float: 'left'}}>文件过滤<span style={{margin: '0 8px 0 2px'}}>:</span></span>
      <Radio.Group
        style={{marginLeft: 20, float: 'left'}}
        value={props.type}
        onChange={this.onChangeFilterRuleType}>
        <Radio value="contain">包含</Radio>
        <Radio value="exclude">排除</Radio>
      </Radio.Group>
      {this.state.full === '1' ? (
        <Button onClick={() => this.setState({full: ''})}>退出全屏</Button>
      ) : (
        <Button onClick={() => this.setState({full: '1'})}>全屏</Button>
      )}
    </div>
  );

  NormalLabel = (props) => (
    <div style={{display: 'inline-block', height: 39, width: 344}}>
      <span style={{float: 'left'}}>{props.title}<span style={{margin: '0 8px 0 2px'}}>:</span></span>
      {this.state.full ? (
        <Button onClick={() => this.setState({full: ''})}>退出全屏</Button>
      ) : (
        <Button onClick={() => this.setStateFull(props.id)}>全屏</Button>
      )}
    </div>
  );

  render() {
    const {info, full, filterRule} = this.state;

    return (
      <React.Fragment>
        <Row>
          <Col span={11}>
            <Form.Item
              colon={false}
              className={full === '1' ? styles.fullScreen : null}
              label={<this.FilterLabel type={filterRule &&  filterRule['type'] || "contain"}/>}
            >
              <Editor
                wrapEnabled
                mode="text"
                theme="tomorrow"
                width="100%"
                height={full === '1' ? '100vh' : '100px'}
                placeholder="每行一条规则"
                value={filterRule && filterRule['data'] || ''}
                onChange={v => filterRule['data'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
            <Form.Item
              colon={false}
              className={full === '3' && styles.fullScreen || null}
              label={<this.NormalLabel title="代码迁出前执行" id="3"/>}>
              <Editor
                mode="sh"
                theme="tomorrow"
                width="100%"
                height={full === '3' && '100vh' || '100px'}
                placeholder="输入要执行的命令"
                value={info['PreCode']}
                onChange={v => info['PreCode'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
            <Form.Item
              colon={false}
              className={full === '5' ? styles.fullScreen : null}
              label={<this.NormalLabel title="应用发布前执行" id="5"/>}>
              <Editor
                mode="sh"
                theme="tomorrow"
                width="100%"
                height={full === '5' ? '100vh' : '100px'}
                placeholder="输入要执行的命令"
                value={info['PreDeploy']}
                onChange={v => info['PreDeploy'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
          </Col>
          <Col span={2}>
            <div className={styles.deployBlock} style={{marginTop: 39}}>
              <Icon type="setting" style={{fontSize: 32}}/>
              <span style={{fontSize: 12, marginTop: 5}}>基础设置</span>
            </div>
            <div className={styles.deployBlock}>
              <Icon type="gitlab" style={{fontSize: 32}}/>
              <span style={{fontSize: 12, marginTop: 5}}>检出代码</span>
            </div>
            <div className={styles.deployBlock}>
              <Icon type="swap" style={{fontSize: 32}}/>
              <span style={{fontSize: 12, marginTop: 5}}>版本切换</span>
            </div>
          </Col>
          <Col span={11}>
            <Form.Item
              colon={false}
              className={full === '2' ? styles.fullScreen : null}
              label={<this.NormalLabel title="自定义全局变量" id="2"/>}>
              <Editor
                wrapEnabled
                mode="text"
                theme="tomorrow"
                width="100%"
                height={full === '2' ? '100vh' : '100px'}
                placeholder="每行一个，例如：HOME=/data/spug"
                value={info['CustomEnvs']}
                onChange={v => info['CustomEnvs'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
            <Form.Item
              colon={false}
              className={full === '4' ? styles.fullScreen : null}
              label={<this.NormalLabel title="代码迁出后执行" id="4"/>}>
              <Editor
                mode="sh"
                theme="tomorrow"
                width="100%"
                height={full === '4' ? '100vh' : '100px'}
                placeholder="输入要执行的命令"
                value={info['PostCode']}
                onChange={v => info['PostCode'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
            <Form.Item
              colon={false}
              className={full === '6' ? styles.fullScreen : null}
              label={<this.NormalLabel title="应用发布后执行" id="6"/>}>
              <Editor
                mode="sh"
                theme="tomorrow"
                width="100%"
                height={full === '6' ? '100vh' : '100px'}
                placeholder="输入要执行的命令"
                value={info['PostDeploy']}
                onChange={v => info['PostDeploy'] = cleanCommand(v)}
                style={{border: '1px solid #e8e8e8'}}/>
            </Form.Item>
          </Col>
        </Row>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button type="primary" onClick={this.handleSubmit}>提交</Button>
          <Button style={{marginLeft: 20}} onClick={this.prePage}>上一步</Button>
        </Form.Item>
      </React.Fragment>
    )
  }
}

export default Ext1Setup3
