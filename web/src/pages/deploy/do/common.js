import React from 'react';
import {connect} from "dva";
import { Steps, Collapse, PageHeader, Spin, Tag, Button, Icon, message} from 'antd';
import history from '@/utils/history';
import {httpGet, httpPost} from '@/utils/request';
import OutView from './OutView';
import styles from './index.module.css';
import lds from 'lodash';
import { stringify } from 'qs';

@connect(({ loading, deploy, config }) => {
  return {
    gitTagList: deploy.gitTagList,
  }
})
class CommonDeploy extends React.Component {
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
        this.setState({
          request: res.data.lists,
          outputs: outputs,
        })
      }).finally(() => this.setState({fetching: false}))
  };

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

        this.socket = new window.WebSocket(`${protocol}//${hostname}:${port}/admin/deploy/ws/${this.id}/ssh/${this.token}`);
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

  // getStatusAlias = () => {
  //   if (Object.keys(this.state.outputs).length !== 0) {
  //     for (let item of [{id: 'local'}, ...this.state.request.Targets]) {
  //       console.log(lds.get(this.state.outputs, `${item.id}.step`))
  //       if (lds.get(this.state.outputs, `${item.id}.status`) === 'error') {
  //         return <Tag color="red">发布异常</Tag>
  //       } else if (lds.get(this.state.outputs, `${item.id}.step`, -1) < 5) {
  //         return <Tag color="blue">发布中</Tag>
  //       }
  //     }
  //     return <Tag color="green">发布成功</Tag>
  //   } else {
  //     return <Tag color="blue">待发布</Tag>
  //   }
  // };

  render() {
    const {AppName, EnvName, Status} = this.state.request;
    return (
      <div>
        <Spin spinning={this.state.fetching}>
          <PageHeader
            title="应用发布"
            subTitle={`服务名：${AppName} —— 环境：${EnvName}`}
            style={{paddingTop: 0}}
            // tags={this.getStatusAlias()}
            extra={this.log ? (
              <Button icon="sync" type="primary" onClick={this.fetch}>刷新</Button>
            ) : (
              <Button icon="play-circle" 
                      type="primary"
                      loading={this.state.loading} 
                      disabled={![2,4].includes(Status)}
                      onClick={this.handleDeploy}>{Status == 2 && "发布" || "回滚"}</Button>
            )}
            onBack={() => history.goBack()}/>

          <Collapse defaultActiveKey={1} className={styles.collapse}>
            <Collapse.Panel showArrow={false} key={1} header={
              <Steps>
                <Steps.Step {...this.getStatus('local', 0)} title="建立连接"/>
                <Steps.Step {...this.getStatus('local', 1)} title="发布准备"/>
                <Steps.Step {...this.getStatus('local', 2)} title="检出前任务"/>
                <Steps.Step {...this.getStatus('local', 3)} title="执行检出"/>
                <Steps.Step {...this.getStatus('local', 4)} title="检出后任务"/>
                <Steps.Step {...this.getStatus('local', 5)} title="执行打包"/>
              </Steps>}>
              <OutView 
                id="local"
                outputs={this.state.outputs}
              />
            </Collapse.Panel>
          </Collapse>

          <Collapse
            defaultActiveKey={'0'}
            className={styles.collapse}
            expandIcon={({isActive}) => <Icon type="caret-right" style={{fontSize: 16}} rotate={isActive ? 90 : 0}/>}>
            { this.state.request.Targets.map((item, index) => (
              <Collapse.Panel key={index} header={
                <div style={{display: 'flex', justifyContent: 'space-between'}}>
                  <b>{item.Title}</b>
                  <Steps size="small" style={{maxWidth: 600}}>
                    <Steps.Step {...this.getStatus(item.ID, 1)} title="数据准备"/>
                    <Steps.Step {...this.getStatus(item.ID, 2)} title="发布前任务"/>
                    <Steps.Step {...this.getStatus(item.ID, 3)} title="执行发布"/>
                    <Steps.Step {...this.getStatus(item.ID, 4)} title="发布后任务"/>
                  </Steps>
                </div>}>
                <OutView 
                  id={item.ID}
                  outputs={this.state.outputs}
                />
              </Collapse.Panel>
            ))}
          </Collapse>
        </Spin>
      </div>
    )
  }
}

export default CommonDeploy

