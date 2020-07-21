import React from 'react';
import { Modal, Steps } from 'antd';
import styles from './index.module.css';
import Setup1 from './Ext2Setup1';
import Setup2 from './Ext2Setup2';
import Setup3 from './Ext2Setup3';


class Ext2Form extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      page: 0,
      info:{},
    }
  }

  componentDidMount() {
    const id=this.props.id
    this.setState({
      info: this.props.info
    })
  };

  handler = (values) => {
    this.setState({
      info: values,
      page: this.state.page + 1,
    })
  }

  prehandler = (values) => {
    this.setState({
      page: this.state.page - 1,
    })
  }

  render() {
    return (
      <Modal
        visible
        width={900}
        maskClosable={false}
        title={Object.keys(this.props.editCacheData).length > 0? '编辑自定义发布' : '新建自定义发布'}
        onCancel={this.props.cancelExt2Visible}
        footer={null}>
        <Steps current={this.state.page} className={styles.steps}>
          <Steps.Step key={0} title="基本配置"/>
          <Steps.Step key={1} title="发布主机"/>
          <Steps.Step key={2} title="执行动作"/>
        </Steps>
        {this.state.page === 0 && 
          <Setup1
            nextPage={this.handler}
            info={this.state.info}
            robotList={this.props.robotList}
          />
        }
        {this.state.page === 1 && 
          <Setup2
            info={this.state.info}
            hostListByAppId={this.props.hostListByAppId}
            nextPage={this.handler}
            prePage={this.prehandler}
            dispatch={this.props.dispatch}
          />
        }
        {this.state.page === 2 && 
          <Setup3
            info={this.state.info}
            nextPage={this.handler}
            prePage={this.prehandler}
            dispatch={this.props.dispatch}
            handleCancel={this.props.cancelExt2Visible}
          />
        }
      </Modal>
    )
  }
}

export default Ext2Form