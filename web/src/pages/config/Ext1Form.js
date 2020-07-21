import React from 'react';
import { Modal, Steps } from 'antd';
import Setup1 from './Ext1Setup1';
import Setup2 from './Ext1Setup2';
import Setup3 from './Ext1Setup3';
import styles from './index.module.css';

class Ext1Form extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      page: 0,
      info:{},
    }
  }

  componentDidMount() {
    const id=this.props.id
    this.setState({
      info: this.props.info,
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
        destroyOnClose= "false"
        title={Object.keys(this.props.editCacheData).length > 0 ? '编辑常规发布' : '新建常规发布'}
        onCancel={this.props.cancelExt1Visible}
        footer={null}>
        <Steps current={this.state.page} className={styles.steps}>
          <Steps.Step key={0} title="基本配置"/>
          <Steps.Step key={1} title="发布主机"/>
          <Steps.Step key={2} title="任务配置"/>
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
            // hostListByAppId={this.props.hostListByAppId}
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
            handleCancel={this.props.cancelExt1Visible}
          />
        }
      </Modal>
    )
  }
}

export default Ext1Form
