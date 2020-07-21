import React from 'react';
import { Link } from 'react-router-dom';
import { Modal, Button, Menu, Spin, Icon, Input, Tooltip } from 'antd';
import styles from './index.module.css';
import lds from 'lodash';

class SelectApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      envId: 0,
      search: '',
      refs: {},
      envIdMap: {},
    }
  }

  componentDidMount() {
    this._initEnv();
  }

  createEnvIdMap = (res) => {
    var tpm = {}
    for (let item of res) {
      this.tpm[item.id] = item
    }

    this.setState({envIdMap: tpm});
  }

  _initEnv = () => {
    if (this.props.configEnvList.length) {
      this.setState({envId: this.props.configEnvList[0].id})
    }
  };

  render() {
    const {envId, envIdMap} = this.state;
    const {configEnvList, projectsList} = this.props;

    let records = projectsList.filter(x => x.EnvId === Number(envId));

    if (this.state.search) {
      records = records.filter(x => x['Name'].toLowerCase().includes(this.state.search.toLowerCase()))
    }
    return (
      <Modal
        visible
        width={800}
        maskClosable={false}
        title="选择应用"
        bodyStyle={{padding: 0}}
        onCancel={this.props.onCancel}
        footer={null}>
        <div className={styles.container}>
          <div className={styles.left}>
            <Spin spinning={configEnvList && false}>
              <Menu
                mode="inline"
                selectedKeys={[String(envId)]}
                style={{border: 'none'}}
                onSelect={({selectedKeys}) => this.setState({envId: selectedKeys[0]})}>
                {configEnvList.map(item => <Menu.Item key={item.id}>{item.Name}</Menu.Item>)}
              </Menu>
            </Spin>
          </div>

          <div className={styles.right}>
            <Spin spinning={projectsList && false}>
              <div className={styles.title}>
                <div>{lds.get(envIdMap, `${envId}.Name`)}</div>
                <Input.Search
                  allowClear
                  style={{width: 200}}
                  placeholder="请输入快速搜应用"
                  onChange={e => this.setState({search: e.target.value})}/>
              </div>
              {records.map(item => (
                <Tooltip key={item.id} >
                  <Button type="primary" className={styles.appBlock} onClick={() => this.props.handleAppChange(item)}>
                    <div 
                         style={{width: 135, overflow: 'hidden', textOverflow: 'ellipsis'}}>
                         {item['Name']}
                    </div>
                  </Button>
                </Tooltip>
              ))}
              {records.length === 0 &&
              <div className={styles.tips}>该环境下还没有可发布的应用哦，快去<Link to="/config/app">应用管理</Link>创建应用发布配置吧。</div>}
            </Spin>
          </div>
        </div>
      </Modal>
    )
  }
}

export default SelectApp
