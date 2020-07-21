import React from 'react';
import { Form, Input, Button, message, Divider, Alert, Icon } from 'antd';
import Editor from 'react-ace';
import 'ace-builds/src-noconflict/mode-sh';
import 'ace-builds/src-noconflict/theme-tomorrow';
import styles from './index.module.css';
import {cleanCommand} from "@/utils/globalTools";

class Ext2Setup3 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      info: {}, 
      PreCode:[],
      PreDeploy: [],
    }
  }
  componentDidMount() {
    var tmpPreCode = this.props.info.hasOwnProperty("PreCode") && 
        this.pareArray(this.props.info['PreCode'].split("|")) || []
    var tmpPreDeploy = this.props.info.hasOwnProperty("PreDeploy") && 
        this.pareArray(this.props.info['PreDeploy'].split("|")) || []
    this.setState({
      info: this.props.info,
      PreCode: tmpPreCode,
      PreDeploy: tmpPreDeploy,
    })
  }

  pareArray = (objectArr) => {
    var arr = []
    objectArr.forEach(item => {
      arr.push(JSON.parse(item))
    })

    return arr
  }

  handleSubmit = () => {
    // this.setState({loading: true});
    var strPreCode = ""
    this.state.PreCode.filter(x => x.title && x.data).forEach(item => {
      if (strPreCode == "") {
        strPreCode = strPreCode + JSON.stringify(item)
      } else {
        strPreCode = strPreCode + "|" + JSON.stringify(item)
      }
    })
    
    var strPreDeploy = ""
    this.state.PreDeploy.filter(x => x.title && x.data).forEach(item => {
      if (strPreDeploy == "") {
        strPreDeploy = strPreDeploy + JSON.stringify(item)
      } else {
        strPreDeploy = strPreDeploy + "|" + JSON.stringify(item)
      }
    })
    
    const info = this.props.info
    info['Extend']    = 2
    info['PreDeploy'] = strPreDeploy
    info['PreCode']   = strPreCode

    if (info.Dtid === undefined) {
      this.props.dispatch({
        type: 'config/deployExtendAdd',
        payload: info,
      }).then(()=> {
        this.props.handleCancel()
      });
    } else {
      this.props.dispatch({
        type: 'config/deployExtendEdit',
        payload: info,
      }).then(()=> {
        this.props.handleCancel()
      }).finally(() => this.setState({loading: false}))
    }
  };

  serverActionsAdd = () => {
    var tmp = this.state.PreCode
    tmp.push({})
    this.setState({
      PreCode: tmp
    })
  }

  serverActionsRemove = (index) => {
    var tmp = this.state.PreCode
    tmp.splice(index, 1)
    this.setState({
      PreCode: tmp
    })
  }

  hostActionsAdd = () => {
    var tmp = this.state.PreDeploy
    tmp.push({})
    this.setState({
      PreDeploy: tmp
    })
  }

  hostActionsRemove = (index) => {
    var tmp = this.state.PreDeploy
    tmp.splice(index, 1)
    this.setState({
      PreDeploy: tmp
    })
  }

  onInputChange = (key, name, index, value) => {
    if (key == "PreCode") {
      var tmp = this.state.PreCode
    } else {
      var tmp = this.state.PreDeploy
    }

    var objtmp = {}
    if(tmp.length > 0 && typeof tmp[index] !== 'undefined') {
      var objtmp = tmp[index]
    }

    objtmp[name] = value
    tmp[index] = objtmp
    this.setState({ 
      info: tmp,
    });
  }

  render() {
    const {PreCode, PreDeploy} = this.state;
    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}} className={styles.ext2Form}>
        {this.state.info.id === undefined && (
          <Alert
            closable
            showIcon
            type="info"
            message="小提示"
            style={{margin: '0 80px 20px'}}
            description={[
              <p key={1}>发布平台 将遵循先本地后目标主机的原则，按照顺序依次执行添加的动作，例如：本地动作1 -> 本地动作2 -> 目标主机动作1 -> 目标主机动作2 ...</p>
            ]}/>
        )}
        {PreCode.map((item, index) => (
          <div key={index} style={{marginBottom: 30, position: 'relative'}}>
            <Form.Item required label={`本地动作${index + 1}`}>
              <Input value={item['title']} onChange={e => this.onInputChange("PreCode", "title", index, e.target.value)} placeholder="请输入"/>
            </Form.Item>

            <Form.Item required label="执行内容">
              <Editor
                wrapEnabled
                mode="sh"
                theme="tomorrow"
                width="100%"
                height="100px"
                value={item['data']}
                onChange={v => this.onInputChange("PreCode", "data", index, cleanCommand(v))}
                placeholder="请输入要执行的动作"/>
            </Form.Item>
            <div className={styles.delAction} onClick={() => this.serverActionsRemove(index)}>
              <Icon type="minus-circle"/>移除
            </div>
          </div>
        ))}
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button type="dashed" block onClick={() => this.serverActionsAdd()}>
            <Icon type="plus"/>添加本地执行动作（在服务端本地执行）
          </Button>
        </Form.Item>
        <Divider/>
        {PreDeploy.map((item, index) => (
          <div key={index} style={{marginBottom: 30, position: 'relative'}}>
            <Form.Item required label={`目标主机动作${index + 1}`}>
              <Input value={item['title']} onChange={e => this.onInputChange("PreDeploy", "title", index, e.target.value)} placeholder="请输入"/>
            </Form.Item>

            <Form.Item required label="执行内容">
              <Editor
                wrapEnabled
                mode="sh"
                theme="tomorrow"
                width="100%"
                height="100px"
                value={item['data']}
                onChange={v => this.onInputChange("PreDeploy", "data", index, cleanCommand(v))}
                placeholder="请输入要执行的动作"/>
            </Form.Item>
            <div className={styles.delAction} onClick={() => this.hostActionsRemove(index)}>
              <Icon type="minus-circle"/>移除
            </div>
          </div>
        ))}
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button type="dashed" block onClick={() => this.hostActionsAdd()}>
            <Icon type="plus"/>添加目标主机执行动作（在部署目标主机执行）
          </Button>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button
            type="primary"
            disabled={[...PreDeploy, ...PreCode].filter(x => x.title && x.data).length === 0}
            loading={this.state.loading}
            onClick={this.handleSubmit}>提交</Button>
          <Button style={{marginLeft: 20}} onClick={() => store.page -= 1}>上一步</Button>
        </Form.Item>
      </Form>
    )
  }
}

export default Ext2Setup3
