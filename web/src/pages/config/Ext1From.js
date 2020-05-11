import React from 'react';
import { Modal, Steps } from 'antd';
import Setup1 from './Ext1Setup1';
import Setup2 from './Ext1Setup2';
import Setup3 from './Ext1Setup3';


class Ext1From extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      page: 0,
      deploy: {},
    }
  }

  componentDidMount() {
    const id=this.props.id
    // this.setState({
    //   page: this.props.page,
    // });
  };

  render() {
    return (
      <Modal
        visible
        width={800}
        maskClosable={false}
        title={Object.keys(this.props.editCacheData).length > 0 ? '编辑常规发布' : '新建常规发布'}
        onCancel={this.props.onCancel}
        footer={null}>
        <Steps current={this.state.page}>
          <Steps.Step key={0} title="基本配置"/>
          <Steps.Step key={1} title="发布主机"/>
          <Steps.Step key={2} title="任务配置"/>
        </Steps>
        {this.state.page === 0 && <Setup1/>}
        {this.state.page === 1 && <Setup2/>}
        {this.state.page === 2 && <Setup3/>}
      </Modal>
    )
  }
}

export default Ext1From
