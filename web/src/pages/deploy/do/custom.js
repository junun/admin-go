import React from 'react';
import {connect} from "dva";
import { Steps, Collapse, PageHeader, Spin, Tag, Button, Icon } from 'antd';
import history from '@/utils/history';
import {httpGet, httpPost} from '@/utils/request';
import OutView from './OutView';;
import styles from './index.module.css';
import lds from 'lodash';
import { stringify } from 'qs';

@connect(({ loading, deploy, config }) => {
  return {
    gitTagList: deploy.gitTagList,
  }
})

class CustomDeploy extends React.Component {
  constructor(props) {
    super(props);
    this.token = sessionStorage.getItem('jwt');
    this.id = props.location.query.id;
    this.log = props.location.query.log === '1' && true || false;
    this.type = props.location.query.type;
    this.templateName = props.location.query.templateName;
    this.state = {
      fetching: false,
      loading: false,
      request: {Targets: []},
      outputs: {},
      PreCode: [],
      PreDeploy: [],
    }
  }

  componentDidMount() {
    this.fetch()
  }

  fetch = () => {
    this.setState({fetching: true});
    var queryParams = {log: this.log}
    httpGet(`/admin/deploy/request/${this.id}?${stringify(queryParams)}`)
      .then(res => {
        const outputs = {}
        while (res.data.lists.Outputs.length) {
          const msg = res.data.lists.Outputs.pop();
          if (!outputs.hasOwnProperty(msg.Key)) {
            const Data = msg.Key === 'local' ? ['读取数据...        '] : [];
            outputs[msg.Key] = {Data}
          }
          this._parse_message(msg, outputs)
        }

        var tmpPreCode = res.data.lists.hasOwnProperty("PreCode") && 
            this.pareArray(res.data.lists['PreCode'].split("|")) || []
        var tmpPreDeploy = res.data.lists.hasOwnProperty("PreDeploy") && 
            this.pareArray(res.data.lists['PreDeploy'].split("|")) || []

        this.setState({
          request: res.data.lists,
          outputs: outputs,
          PreCode: tmpPreCode,
          PreDeploy: tmpPreDeploy,
        })
      }).finally(() => this.setState({fetching: false}))
  };

  pareArray = (objectArr) => {
    var arr = []
    objectArr.forEach(item => {
      arr.push(JSON.parse(item))
    })

    return arr
  }

  _parse_message = (message, outputs) => {
    outputs = outputs || this.state.outputs;
    const {Key, Data, Step, Status} = message;
    if (Data !== undefined) {
      outputs[Key]['Data'].push(Data);
    }

    if (Step !== undefined) outputs[Key]['step'] = Step;
    if (Status !== undefined) outputs[Key]['status'] = Status;

    this.setState({
      outputs: outputs,
    })
  };

  handleDeploy = () => {
    this.setState({loading: true});
    httpPost(`/admin/deploy/request/${this.id}`)
      .then(res => {
        if (res.code != 200) { 
          message.error(res.message)
          return
        }

        var tmp = this.state.request
        tmp.Status = '2'
        this.setState({
          request: tmp,
          outputs: res.data,
        })

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const hostname = window.location.hostname;
        const port = window.location.port;

        this.socket = new window.WebSocket(`${protocol}//${hostname}:${port}/admin/deploy/ws/${id}/ssh/${token}`);
        // this.socket = new window.WebSocket(`${protocol}//127.0.0.1:9090/admin/deploy/ws/${this.id}/ssh/${this.token}`);

        this.socket.onopen = () => {
          this.socket.send('ok');
        };
        this.socket.onmessage = e => {
          if (e.data === 'pong') {
            this.socket.send('ping')
          } else {
            this._parse_message(JSON.parse(e.data))
          }
        }
      })
      .finally(() => this.setState({loading: false}))
  };

  getStatus = (key, n) => {
    const step = lds.get(this.state.outputs, `${key}.step`, -1);

    const isError = lds.get(this.state.outputs, `${key}.status`) === 'error';
    const icon = <Icon type="loading"/>;
    if (n > step) {
      return {key: n, status: 'wait'}
    } else if (n === step) {
      return isError ? {key: n, status: 'error'} : {key: n, status: 'process', icon}
    } else {
      return {key: n, status: 'finish'}
    }
  };

  render() {
    const {AppName, EnvName, Status} = this.state.request;
    const {PreCode, PreDeploy} = this.state;

    return (
      <div>
        <Spin spinning={this.state.fetching}>
          <PageHeader
            title="应用发布"
            subTitle={`服务名：${AppName} —— 环境：${EnvName}`}
            style={{padding: 0}}
            // tags={this.getStatusAlias()}
            extra={this.log ? (
              <Button icon="sync" type="primary" onClick={this.fetch}>刷新</Button>
            ) : (
              <Button icon="play-circle" loading={this.state.loading} type="primary"
                      loading={this.state.loading} 
                      disabled={![2,4].includes(Status)}
                      onClick={this.handleDeploy}>{Status == 2 && "发布" || "回滚"}</Button>
            )}
            onBack={() => history.goBack()}/>
          <Collapse defaultActiveKey={1} className={styles.collapse}>
            <Collapse.Panel showArrow={false} key={1} header={
              <Steps style={{maxWidth: 400 + PreCode.length * 200}}>
                <Steps.Step {...this.getStatus('local', 0)} title="建立连接"/>
                <Steps.Step {...this.getStatus('local', 1)} title="发布准备"/>
                {PreCode.map((item, index) => (
                  <Steps.Step {...this.getStatus('local', 2 + index)} key={index} title={item.title}/>
                ))}
              </Steps>}>
              <OutView 
                id="local"
                outputs={this.state.outputs}
              />
            </Collapse.Panel>
          </Collapse>

          {PreDeploy.length > 0 && (
            <Collapse
              defaultActiveKey={'0'}
              className={styles.collapse}
              expandIcon={({isActive}) => <Icon type="caret-right" style={{fontSize: 16}} rotate={isActive ? 90 : 0}/>}>
              { this.state.request.Targets.map((item, index) => (
                <Collapse.Panel key={index} header={
                  <div style={{display: 'flex', justifyContent: 'space-between'}}>
                    <b>{item.Title}</b>
                    <Steps size="small" style={{maxWidth: 150 + PreDeploy.length * 150}}>
                      <Steps.Step {...this.getStatus(item.ID, 1)} title="数据准备"/>
                      {PreDeploy.map((action, index) => (
                        <Steps.Step {...this.getStatus(item.ID, 2 + index)} key={index} title={action.title}/>
                      ))}
                    </Steps>
                  </div>}>
                  <OutView 
                    id={item.ID}
                    outputs={this.state.outputs}
                  />
                </Collapse.Panel>
              ))}
            </Collapse>
          )}
          {PreDeploy.length === 0 && this.state.fetching === false && (
            <div className={styles.ext2Tips}>无目标主机动作</div>
          )}
        </Spin>
      </div>
    )
  }
}

export default CustomDeploy
